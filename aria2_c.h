#ifndef ARIA2_C_H
#define ARIA2_C_H

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Type definition for torrent information.
 */
struct TorrentInfo {
  int totalFile;
  const char *infoHash;
  struct FileInfo *files;
};

/**
 * Type definition for file information in torrent.
 */
struct FileInfo {
  int index;
  const char *name;
  int64_t length;
  bool selected;
};

/**
 * Type definition for download information.
 */
struct DownloadInfo {
  int status;
  int64_t totalLength;
  int64_t bytesCompleted;
  int64_t uploadLength;
  int downloadSpeed;
  int uploadSpeed;
};

int init(uint64_t aria2goPointer);
uint64_t addUri(char *uri);
struct TorrentInfo *parseTorrent(char *fp);
uint64_t addTorrent(char *fp, const char *options);
bool changeOptions(uint64_t gid, const char *options);
void start();
bool pause(uint64_t gid);
bool resume(uint64_t gid);
bool removeDownload(uint64_t gid);
struct DownloadInfo *getDownloadInfo(uint64_t gid);

#ifdef __cplusplus
}
#endif

#endif