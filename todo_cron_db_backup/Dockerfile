FROM busybox:1.38

COPY scripts/db-backup-task.sh /usr/local/bin/db-backup-task.sh

RUN chmod +x /usr/local/bin/db-backup-task.sh

CMD [ "/usr/local/bin/db-backup-task.sh" ]
