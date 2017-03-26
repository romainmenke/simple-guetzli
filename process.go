package main

func needsProc(j *job) bool {

	var imgChanged bool
	var settingsChanged bool

	if j.report.Empty() {
		settingsChanged = true
		imgChanged = true
	}

	// if !j.report.Empty() && j.report.Version != j.settings.version {
	// 	settingsChanged = true
	// }
	//
	// if !j.report.Empty() && j.report.Quality != j.settings.quality {
	// 	settingsChanged = true
	// }

	modTime := timeModified(j.settings.source + j.fileName)
	if !j.report.Empty() && !modTime.Equal(j.report.ModTime) {
		imgChanged = true
	}

	sha := sha1ForFile(j.settings.source + j.fileName)
	if !j.report.Empty() && sha != j.report.Sha1 {
		imgChanged = true
	}

	if !imgChanged && !settingsChanged {
		return false
	}

	j.report.Quality = j.settings.quality
	j.report.ModTime = modTime
	j.report.Path = j.settings.source + j.fileName
	j.report.Sha1 = sha
	j.report.Version = j.settings.version

	return true

}
