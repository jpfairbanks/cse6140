/* hashes.go: utilities to help with hashing things.
Author: James Fairbanks
Liscence: BSD

 we need to find primes and compute a*x+b mod p as the hash,
then we need to test that these hashes are indeed pairwise independent for
distinct values of a and b.
*/

package hashes

import(
	"fmt"
)


