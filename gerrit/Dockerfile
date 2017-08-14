FROM concourse/buildroot:git

WORKDIR /opt/resource

ADD gerrit-resource .
RUN chmod +x gerrit-resource
RUN ln -s gerrit-resource check
RUN ln -s gerrit-resource in
RUN ln -s gerrit-resource out
