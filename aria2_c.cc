#include "aria2_c.h"
#include "_cgo_export.h"
#include <aria2/aria2.h>
#include <iostream>
#include <sstream>
#include <string.h>

std::vector<std::string> splitBySemicolon(std::string in) {
  std::string val;
  std::stringstream is(in);
  std::vector<std::string> out;

  while (getline(is, val, ';')) {
    out.push_back(val);
  }
  return out;
}

const char *toCStr(std::string in) {
  int len = in.length();
  char *val = new char[len + 1];
  memcpy(val, in.data(), len);
  val[len] = '\0';
  return val;
}

const char *getFileName(std::string dir, std::string path) {
  if (path.find(dir, 0) == std::string::npos) {
    return toCStr(path);
  }
  int index = dir.size();
  std::string name = path.substr(index + 1);
  return toCStr(name);
}

const char *toAria2goOptions(aria2::KeyVals options) {
  std::vector<std::pair<std::string, std::string>>::iterator it;

  std::string cOptions;
  for (it = options.begin(); it != options.end(); it++) {
    std::pair<std::string, std::string> p = *it;
    cOptions += p.first + ";" + p.second + ";";
  }

  return toCStr(cOptions);
}

aria2::KeyVals toAria2Options(const char *options) {
  aria2::KeyVals aria2Options;

  if (options == nullptr) {
    return aria2Options;
  }

  std::vector<std::string> o = splitBySemicolon(std::string(options));
  /* key and val should be pair */
  if (o.size() % 2 != 0) {
    return aria2Options;
  }

  for (int i = 0; i < o.size(); i += 2) {
    std::string key = o[i];
    std::string val = o[i + 1];
    aria2Options.push_back(std::make_pair(key, val));
  }

  return aria2Options;
}

struct FileInfo *parseFileData(aria2::DownloadHandle *dh) {
  std::string dir = dh->getDir();
  std::vector<aria2::FileData> files = dh->getFiles();
  int numFiles = dh->getNumFiles();
  struct FileInfo *allFiles = new FileInfo[numFiles];
  for (int i = 0; i < numFiles; i++) {
    aria2::FileData file = files[i];
    struct FileInfo *fi = new FileInfo();
    fi->index = file.index;
    fi->name = getFileName(dir, file.path);
    fi->length = file.length;
    fi->completedLength = file.completedLength;
    fi->selected = file.selected;

    allFiles[i] = *fi;
  }
  return allFiles;
}

/**
 * Global aria2 session.
 */
aria2::Session *session;
/**
 * Global aria2go go pointer.
 */
uint64_t aria2goPointer;

/**
 * Download event callback for aria2.
 */
int downloadEventCallback(aria2::Session *session, aria2::DownloadEvent event,
                          const aria2::A2Gid gid, void *userData) {
  notifyEvent(aria2goPointer, gid, event);
  return 0;
}

/**
 * Initial aria2 library.
 */
int init(uint64_t pointer, const char *options) {
  aria2goPointer = pointer;
  int ret = aria2::libraryInit();
  aria2::SessionConfig config;
  config.keepRunning = true;
  /* do not use signal handler, cause c will block go */
  config.useSignalHandler = false;
  config.downloadEventCallback = downloadEventCallback;
  session = aria2::sessionNew(toAria2Options(options), config);
  return ret;
}

/**
 * Deinit aria2 library, this must be invoked when process exit(signal handler
 * is not used), so aria2 will be able to save session config.
 */
int deinit() {
  int ret = aria2::sessionFinal(session);
  aria2::libraryDeinit();
  return ret;
}

/**
 * Adds new HTTP(S)/FTP/BitTorrent Magnet URI. See `addUri` in aria2.
 *
 * @param uri uri to add
 */
uint64_t addUri(char *uri, const char *options) {
  std::vector<std::string> uris = {uri};
  aria2::A2Gid gid;
  int ret = aria2::addUri(session, &gid, uris, toAria2Options(options));
  if (ret < 0) {
    return 0;
  }

  return gid;
}

/**
 * Parse torrent file for given file path. This will just retrieve related
 * information from torrent file, and aria2 will not download it.
 */
struct TorrentInfo *parseTorrent(char *fp) {
  aria2::KeyVals options;
  options.push_back(std::make_pair("dry-run", "true"));

  aria2::A2Gid gid;
  int ret = aria2::addTorrent(session, &gid, fp, options);
  if (ret < 0) {
    return nullptr;
  }

  aria2::DownloadHandle *dh = aria2::getDownloadHandle(session, gid);
  if (!dh) {
    return nullptr;
  }

  std::vector<aria2::FileData> files = dh->getFiles();
  aria2::BtMetaInfoData btMetaInfo = dh->getBtMetaInfo();

  struct TorrentInfo *ti = new TorrentInfo();
  int numFiles = files.size();
  ti->numFiles = dh->getNumFiles();
  ti->infoHash = toCStr(dh->getInfoHash());

  /* retrieve all BitTorrent meta information */
  struct MetaInfo *mi = new MetaInfo();
  mi->name = toCStr(btMetaInfo.name);
  mi->comment = toCStr(btMetaInfo.comment);
  mi->creationUnix = btMetaInfo.creationDate;
  std::vector<std::vector<std::string>> announceList = btMetaInfo.announceList;
  std::vector<std::vector<std::string>>::iterator it;
  std::string cAnnounceList;
  for (it = announceList.begin(); it != announceList.end(); it++) {
    std::vector<std::string>::iterator cit;
    std::vector<std::string> childList = *it;
    for (cit = childList.begin(); cit != childList.end(); cit++) {
      cAnnounceList += *cit;
      if (it != announceList.end() - 1 || cit != childList.end() - 1) {
        cAnnounceList += ";";
      }
    }
  }
  mi->announceList = toCStr(cAnnounceList);
  ti->metaInfo = mi;
  /* retrieve all files information */
  ti->files = parseFileData(dh);

  /* remove from download queue, just retrieve information */
  aria2::deleteDownloadHandle(dh);
  aria2::removeDownload(session, gid);

  return ti;
}

/**
 * Add bit torrent file. See `addTorrent` in aria2.
 */
uint64_t addTorrent(char *fp, const char *options) {
  aria2::A2Gid gid;

  int ret = aria2::addTorrent(session, &gid, fp, toAria2Options(options));
  if (ret < 0) {
    return 0;
  }

  return gid;
}

/**
 * Change aria2 options. See `changeOption` in aria2.
 */
bool changeOptions(uint64_t gid, const char *options) {
  return aria2::changeOption(session, gid, toAria2Options(options)) == 0;
}

/**
 * Get options for given gid. see `getOptions` in aria2.
 */
const char *getOptions(uint64_t gid) {
  aria2::DownloadHandle *dh = aria2::getDownloadHandle(session, gid);
  if (!dh) {
    return nullptr;
  }

  return toAria2goOptions(dh->getOptions());
}

/**
 * Change global options. See `changeGlobalOption` in aria2.
 */
bool changeGlobalOptions(const char *options) {
  return aria2::changeGlobalOption(session, toAria2Options(options));
}

/**
 * Get global options. see `getGlobalOptions` in aria2.
 */
const char *getGlobalOptions() {
  aria2::KeyVals options = aria2::getGlobalOptions(session);
  return toAria2goOptions(options);
}

/**
 * Start to download. This will block current thread.
 */
void start() {
  for (;;) {
    int ret = aria2::run(session, aria2::RUN_DEFAULT);
    if (ret != 1) {
      break;
    }
  }
}

/**
 * Pause an active download with given gid. This will mark the download to
 * `DOWNLOAD_PAUSED`. See `resume`.
 */
bool pause(uint64_t gid) { return aria2::pauseDownload(session, gid) == 0; }

/**
 * Resume a paused download with given gid. See `pause`.
 */
bool resume(uint64_t gid) { return aria2::unpauseDownload(session, gid) == 0; }

/**
 * Remove a download in queue with given gid. This will stop downloading and
 * seeding(for torrent).
 */
bool removeDownload(uint64_t gid) {
  return aria2::removeDownload(session, gid) == 0;
}

/**
 * Get download information for current download with given gid.
 */
struct DownloadInfo *getDownloadInfo(uint64_t gid) {
  aria2::DownloadHandle *dh = aria2::getDownloadHandle(session, gid);
  if (!dh) {
    return nullptr;
  }

  struct DownloadInfo *di = new DownloadInfo();
  di->status = dh->getStatus();
  di->totalLength = dh->getTotalLength();
  di->bytesCompleted = dh->getCompletedLength();
  di->uploadLength = dh->getUploadLength();
  di->downloadSpeed = dh->getDownloadSpeed();
  di->uploadSpeed = dh->getUploadSpeed();
  di->pieceLength = dh->getPieceLength();
  di->numPieces = dh->getNumPieces();
  di->connections = dh->getConnections();
  di->numFiles = dh->getNumFiles();
  di->files = parseFileData(dh);

  /* delete download handle */
  aria2::deleteDownloadHandle(dh);

  return di;
}
