# VideoSender

Для запуска нужно заполнить следующие переменные

```bash
export VS_TG_TOKEN="change_me"
export VS_TG_GROUP="change_me"
export VS_REDIS_HOST=localhost:6379
export VS_REDIS_PASSWORD="change_me"
```

Далее послать POST запрос

```bash
    curl -XPOST http://localhost:8090/addjob -d '{
        "key": "1713342768:camera2",
        "ttl": 0,
        "value": {
            "video_file": "/mnt/video_104-79-20251021-100524.mp4",
            "camera_name": "camera2",
            "file_size":  15495515
        }
    }'
```
