FROM scratch

COPY configs/api.config.yaml.dist configs/api.config.yaml
COPY beanstalk-api /

ENTRYPOINT ["/beanstalk-api"]
