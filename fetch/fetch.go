package fetch

func (h Handle) URL(pkgbase string) string {
	return h.AURURL + "/" + pkgbase + ".git"
}

func (h Handle) Download(pkgs ...string) error {
	for _, pkg := range pkgs {
		_, err := h.GitHandle.Download(h.URL(pkg), h.CacheDir, pkg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h Handle) Merge(pkgs ...string) error {
	for _, pkg := range pkgs {
		needMerge, err := h.GitHandle.NeedMerge(h.CacheDir, pkg)
		if err != nil {
			return err
		}

		if !needMerge {
			continue
		}

		err = h.GitHandle.Merge(h.CacheDir, pkg)
		if err != nil {
			return err
		}
	}

	return nil
}
