youtrack:
	docker run -it --name youtrack \
	-v ~/youtrack/data:/opt/youtrack/data \
	-v ~/youtrack/conf:/opt/youtrack/conf \
	-v ~/youtrack/logs:/opt/youtrack/logs \
	-v ~/youtrack/backups:/opt/youtrack/backups \
	-p 8080:8080 \
	jetbrains/youtrack:2025.1.74704
