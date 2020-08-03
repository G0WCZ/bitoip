/*
Copyright (C) 2020 Graeme Sutherland, Nodestone Limited


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
#include <inttypes.h>
#include <sys/time.h>

/**
 * A time64 64 bit timestamp for now
 */
uint64_t timestamp_64_now() {
    struct timeval tv_now;

    gettimeofday(&tv_now, NULL);
    return (uint64_t)(tv_now.tv_sec * 1000000L + (uint64_t)tv_now.tv_usec) * 1000L;
}

