#ifndef __sys_info__
#define __sys_info__
struct sysinfo {
  uint64 freemem;   // amount of free memory (bytes)
  uint64 nproc;     // number of process
};
#endif
