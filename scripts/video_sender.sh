#!/bin/bash

sendVideo() {
	CAMERA=$1
	VIDEO_FILE=$2
	DATE=$(date +%s)
	FILE_SIZE=$(stat -c %s "$VIDEO_FILE")

	curl -XPOST http://localhost:8090/addjob  -d "{
	\"key\": \"$DATE:$CAMERA\",
	\"ttl\": 0,
	\"value\":
		{
			\"video_file\": \"$VIDEO_FILE\",
			\"camera_name\": \"$CAMERA\",
			\"file_size\": $FILE_SIZE
		}
	}"
}

sendVideo "$@"
