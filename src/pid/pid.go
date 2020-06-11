
package pid

import (
	"os"
	"strconv"
)

func Get() (int) {
	return os.Getpid()
}

func Save(path string) (int, error) {

	pid := Get()
	fd, err := os.Create(path)
	if err != nil {
    	return -1, err
	}

	_, err = fd.WriteString(strconv.Itoa(pid))

	if err != nil {
		return -1, err
	}

	return pid, nil
}
