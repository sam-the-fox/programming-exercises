local fox FORCED to write javascript, 12 dead 580 injured

these are exercises and not real projects, so i license them under the Unlicense.
see UNLICENSE in this repository for details, or <https://unlicense.org/>.

## projects

### A004

silly little status monitor in both python (with flask + requests) and go (only the standard library). both use data from `A004/services.json`.
the go version has such creature comforts as "heeds an interrupt signal" and "is concurrency safe". the python version has neither, and abuses `threading.Thread` because i couldn't use celery.
