module bitbucket.org/HeilaSystems/session

go 1.14

replace (
	bitbucket.org/HeilaSystems/dependencybundler v0.0.2 => ../dependencybundler
	bitbucket.org/HeilaSystems/log v0.0.0 => ../log
	bitbucket.org/HeilaSystems/session v0.0.0 => ./
	bitbucket.org/HeilaSystems/trace v0.0.0 => ../trace
	bitbucket.org/HeilaSystems/transport v0.0.0 => ../transport
)

require bitbucket.org/HeilaSystems/dependencybundler v0.0.2
