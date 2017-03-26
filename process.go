package main

func needsProc(j *job) bool {

	var imgChanged bool
	var settingsChanged bool

	if j.report == nil {
		settingsChanged = true
		imgChanged = true
	}

	// if j.report != nil && j.report.Version != j.settings.version {
	// 	settingsChanged = true
	// }
	//
	// if j.report != nil && j.report.Quality != j.settings.quality {
	// 	settingsChanged = true
	// }

	modTime := timeModified(j.settings.source + j.fileName)
	if j.report != nil && !modTime.Equal(j.report.ModTime) {
		imgChanged = true
	}

	sha := sha1ForFile(j.settings.source + j.fileName)
	if j.report != nil && sha != j.report.Sha1 {
		imgChanged = true
	}

	if !imgChanged && !settingsChanged {
		return false
	}

	j.report = &guetzliReport{
		Quality: j.settings.quality,
		ModTime: modTime,
		Path:    j.settings.source + j.fileName,
		Sha1:    sha,
		Version: j.settings.version,
	}

	return true

}
