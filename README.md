cse6140
=======

My programming assignments for class CSE6140 Algorithms

Each program can be found in its own subdirectory. These will be given useful names because then the Go import paths are meaningful. I intend to do all of the coding for this class in Go (golang). This README file should contain a manifest mapping the assignments to package names.

I am structuring things this way so that I can use github for managing the code but not make a separate repo for each project because it seems wasteful to make a repo that only has active commits for a 1-2 week period. 

Feel free to use this code for any purpose, academic, instructional, commercial. Please let me know that you are using it. I am curious about my webtraffic. If you are using this code in an academic paper or homework assignment, then you must cite me in order to avoid plagarism.

- [X] Assignment 1 : strides : strided array access to test cache performance.
- [ ] Project : Parallel Count Min Sketch


TODOs
=====
- [ ] Batch insertions
- [ ] Improve cache performance by sorting the insertions for each row.
- [ ] Element level parallelism
- [ ] Signal based parallel insert for Depth > 1000
- [ ] Determine good parameters
- [ ] Disable exact answer computation
- [ ] More Tests
- [ ] Compute speedup against go's builting map[int64][int64] type
- [ ] Compute parallel speedup for reasonable parameters
