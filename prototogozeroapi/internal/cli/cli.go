package cli

import (
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"

	"prototogozeroapi/internal/genapi"
	"prototogozeroapi/internal/protosrc"
)

type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func Run(args []string) error {
	fs := flag.NewFlagSet("prototogozeroapi", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var includes multiFlag
	fs.Var(&includes, "I", "proto import include path")

	if err := fs.Parse(args); err != nil {
		return err
	}
	rest := fs.Args()
	if len(rest) < 1 {
		return errors.New("missing proto file path")
	}

	in := rest[0]
	dir := filepath.Dir(in)
	base := strings.TrimSuffix(filepath.Base(in), filepath.Ext(in))
	outPath := filepath.Join(dir, base+".api")

	loaded, err := protosrc.Load(in, []string(includes))
	if err != nil {
		return err
	}
	api, err := genapi.Generate(loaded)
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(api), 0o644)
}

