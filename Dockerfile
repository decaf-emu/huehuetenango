FROM golang:1.8

# switch the shell used by RUN from sh to bash
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

ENV OUTPUT_DIR /opt/huehuetenango
ENV DATA_DIR /data/huehuetenango

# create the huehuetenango user
RUN mkdir -p $OUTPUT_DIR && \
  mkdir -p $DATA_DIR && \
  groupadd -r huehuetenango && \
  useradd -r -u 528 -g huehuetenango -d $OUTPUT_DIR -s /sbin/nologin huehuetenango && \
  chown -R 528:huehuetenango $OUTPUT_DIR $DATA_DIR && \
  chmod 775 $OUTPUT_DIR $DATA_DIR

# create and switch to the huehuetenango user
VOLUME /data

ENV HOME /root
ENV GOPATH $HOME/go
ENV PROJECT_DIR $GOPATH/src/github.com/decaf-emu/huehuetenango

# copy the project files
RUN mkdir -p $PROJECT_DIR
COPY . $PROJECT_DIR
WORKDIR /root/go/src/github.com/decaf-emu/huehuetenango

# install nvm
RUN curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.1/install.sh | bash

# install node
ENV NVM_DIR $HOME/.nvm
ENV NPM_DIR $HOME/.npm
ENV YARN_DIR $HOME/.yarn
ENV PATH $YARN_DIR/bin:$GOPATH/bin:$PATH

RUN source $NVM_DIR/nvm.sh && \
  cd static && \
  nvm install

# install yarn
RUN source $NVM_DIR/nvm.sh && \
  curl -o- -L https://yarnpkg.com/install.sh | bash

# build huehuetenango
RUN source $NVM_DIR/nvm.sh && \
  cd static && \
  nvm use && \
  cd .. && \
  make

# copy and assign permissions to the build files
RUN mkdir -p $OUTPUT_DIR && \
  cp $PROJECT_DIR/huehuetenango $OUTPUT_DIR/huehuetenango && \
  cp -R $PROJECT_DIR/static/dist $OUTPUT_DIR/static && \
  chown -R huehuetenango:huehuetenango $OUTPUT_DIR

# allow huehuetenango to use ports less than 1024
RUN setcap 'cap_net_bind_service=+ep' $OUTPUT_DIR/huehuetenango

# clean the build related directories
RUN rm -rf $GOPATH && \
  rm -rf $NVM_DIR && \
  rm -rf $YARN_DIR && \
  rm -rf $NPM_DIR

# install and configure supervisor
RUN apt-get update && \
  apt-get install -y supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

EXPOSE 80 443

USER huehuetenango
WORKDIR $OUTPUT_DIR
CMD ["/usr/bin/supervisord"]
