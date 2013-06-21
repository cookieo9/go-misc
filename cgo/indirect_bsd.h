// +build darwin freebsd netbsd openbsd

#include <stdio.h>

extern int (*c_reader)(void *, char *, int);
extern int (*c_writer)(void *, const char *, int);
extern fpos_t (*c_seeker)(void *, fpos_t, int);
extern int (*c_closer)(void *);
