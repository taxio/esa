package main

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/srvc/fail/v4"
)

var digitReg = regexp.MustCompile(`^-?\d+$`)
var ErrInvalidPostIdFormat = errors.New("invalid post id format")

func ParsePostIdFromArg(arg string) (int, error) {
	var pIdStr string
	if digitReg.Match([]byte(arg)) {
		pIdStr = arg
	} else {
		u, err := url.Parse(arg)
		if err != nil {
			return 0, fail.Wrap(err)
		}
		elms := strings.Split(u.Path, "/")
		pIdStr = elms[len(elms)-1]
	}

	pId, err := strconv.Atoi(pIdStr)
	if err != nil {
		return 0, fail.Wrap(err)
	}

	if pId <= 0 {
		return 0, fail.Wrap(ErrInvalidPostIdFormat)
	}

	return pId, nil
}
