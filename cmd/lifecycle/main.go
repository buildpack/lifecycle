package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"

	"github.com/buildpacks/lifecycle"
	"github.com/buildpacks/lifecycle/buildpack"
	"github.com/buildpacks/lifecycle/cache"
	"github.com/buildpacks/lifecycle/cmd"
)

func main() {
	platformAPI := cmd.EnvOrDefault(cmd.EnvPlatformAPI, cmd.DefaultPlatformAPI)
	if err := cmd.VerifyPlatformAPI(platformAPI); err != nil {
		cmd.Exit(err)
	}

	switch strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])) {
	case "detector":
		cmd.Run(&detectCmd{detectArgs: detectArgs{platformAPI: platformAPI}}, false)
	case "analyzer":
		cmd.Run(&analyzeCmd{analyzeArgs: analyzeArgs{platformAPI: platformAPI}}, false)
	case "restorer":
		cmd.Run(&restoreCmd{platformAPI: platformAPI}, false)
	case "builder":
		cmd.Run(&buildCmd{buildArgs: buildArgs{platformAPI: platformAPI}}, false)
	case "exporter":
		cmd.Run(&exportCmd{exportArgs: exportArgs{platformAPI: platformAPI}}, false)
	case "rebaser":
		cmd.Run(&rebaseCmd{platformAPI: platformAPI}, false)
	case "creator":
		cmd.Run(&createCmd{platformAPI: platformAPI}, false)
	default:
		if len(os.Args) < 2 {
			cmd.Exit(cmd.FailCode(cmd.CodeInvalidArgs, "parse arguments"))
		}
		if os.Args[1] == "-version" {
			cmd.ExitWithVersion()
		}
		subcommand()
	}
}

func subcommand() {
	phase := filepath.Base(os.Args[1])
	switch phase {
	case "detect":
		cmd.Run(&detectCmd{}, true)
	case "analyze":
		cmd.Run(&analyzeCmd{}, true)
	case "restore":
		cmd.Run(&restoreCmd{}, true)
	case "build":
		cmd.Run(&buildCmd{}, true)
	case "export":
		cmd.Run(&exportCmd{}, true)
	case "rebase":
		cmd.Run(&rebaseCmd{}, true)
	case "create":
		cmd.Run(&createCmd{}, true)
	default:
		cmd.Exit(cmd.FailCode(cmd.CodeInvalidArgs, "unknown phase:", phase))
	}
}

func verifyBuildpackApis(group buildpack.Group) error {
	for _, bp := range group.Group {
		if bp.API == "" {
			// if this group was generated by this lifecycle bp.API should be set
			// but if for some reason it isn't default to 0.2
			bp.API = "0.2"
		}
		if err := cmd.VerifyBuildpackAPI(bp.String(), bp.API); err != nil {
			return err
		}
	}
	return nil
}

func initCache(cacheImageTag, cacheDir string, keychain authn.Keychain) (lifecycle.Cache, error) {
	var (
		cacheStore lifecycle.Cache
		err        error
	)
	if cacheImageTag != "" {
		cacheStore, err = cache.NewImageCacheFromName(cacheImageTag, keychain)
		if err != nil {
			return nil, cmd.FailErr(err, "create image cache")
		}
	} else if cacheDir != "" {
		cacheStore, err = cache.NewVolumeCache(cacheDir)
		if err != nil {
			return nil, cmd.FailErr(err, "create volume cache")
		}
	}
	return cacheStore, nil
}
