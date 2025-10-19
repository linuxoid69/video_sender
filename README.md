# VideoSender

Для запуска нужно заполнить следующие переменные

```bash
export TG_TOKEN=123
export TG_GROUP=123
```

Далее послать POST запрос

```bash
curl -XPOST http://localhost:8090/video  -d '{"video_path":"file_name", "camera":"camera_name"}"'
'
```
