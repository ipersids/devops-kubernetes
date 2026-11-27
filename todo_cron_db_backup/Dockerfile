FROM postgres:18-alpine

ENV CLOUDSDK_PYTHON=/usr/bin/python3
ENV PATH="/usr/local/gcloud/google-cloud-sdk/bin:${PATH}"

COPY scripts/db-backup-task.sh /usr/local/bin/db-backup-task.sh

RUN apk add --no-cache curl python3 ca-certificates

WORKDIR /usr/local/gcloud

RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-linux-x86_64.tar.gz \
  && tar -xf google-cloud-cli-linux-x86_64.tar.gz \
  && ./google-cloud-sdk/install.sh --quiet --usage-reporting=false --path-update=false \
  && rm google-cloud-cli-linux-x86_64.tar.gz \
  && rm -rf /root/.config/gcloud

RUN chmod +x /usr/local/bin/db-backup-task.sh

CMD [ "/usr/local/bin/db-backup-task.sh" ]
