package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const defaultOut = "./moonguard-clients"

type langConfig struct {
	execFlag string
}

var supportedLangs = map[string]langConfig{
	"go": {
		execFlag: "--go_out=plugins=grpc:",
	},
	"cpp": {
		execFlag: "--cpp_out=",
	},
	"csharp": {
		execFlag: "--csharp_out=",
	},
	"java": {
		execFlag: "--java_out=",
	},
	"js": {
		execFlag: "--js_out=",
	},
	"objc": {
		execFlag: "--objc_out=",
	},
	"php": {
		execFlag: "--php_out=",
	},
	"python": {
		execFlag: "--python_out=",
	},
	"ruby": {
		execFlag: "--ruby_out=",
	},
}

func buildGrpcCommand(sources []string, langs []string, outDir string) (*exec.Cmd, error) {
	grpcExec, err := exec.LookPath("protoc")
	if err != nil {
		return nil, fmt.Errorf("unable to find command `protoc` in path")
	}
	args := []string{grpcExec}

	for _, lang := range langs {
		cfg := supportedLangs[lang]
		langOutdir := path.Join(outDir, lang)

		os.MkdirAll(langOutdir, 0777)

		flag := cfg.execFlag + langOutdir
		args = append(args, flag)
	}

	args = append(args, sources...)

	cmd := &exec.Cmd{
		Path:   grpcExec,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Args:   args,
	}

	fmt.Printf("%v", args)

	return cmd, nil
}

func findInputSources(pattern string) ([]string, error) {
	if pattern == "" {
		return nil, fmt.Errorf("first argument must be your protobuf source path")
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("unable to find input sources: %s", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("found no protobuf sources matching `%s`", matches)
	}

	return matches, nil
}

func validateLangs(langs []string) error {
	for _, lang := range langs {
		if _, ok := supportedLangs[lang]; ok != true {
			return fmt.Errorf("`%s` is not yet supported by the moonguard client generator", lang)
		}
	}
	return nil
}

func genAction(c *cli.Context) error {
	inputSourcesStr := c.Args().Get(0)
	sources, err := findInputSources(inputSourcesStr)
	if err != nil {
		return err
	}

	outdir := c.String("out")
	if outdir == "" {
		outdir = defaultOut
	}

	langs := c.StringSlice("languages")
	err = validateLangs(langs)
	if err != nil {
		return err
	}

	cmd, err := buildGrpcCommand(sources, langs, outdir)
	if err != nil {
		return err
	}

	return cmd.Run()
}

func GetGenCommand() *cli.Command {
	return &cli.Command{
		Name:   "gen",
		Action: genAction,
		Usage:  "generate moonguard gRPC clients",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "out",
				Usage:       "Output directory for generated clients",
				DefaultText: defaultOut,
			},
			&cli.StringSliceFlag{
				Name:     "languages",
				Aliases:  []string{"l"},
				Usage:    "Build gRPC clients for this set of languages",
				Required: true,
			},
		},
	}
}
