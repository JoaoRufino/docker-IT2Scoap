FROM pritunl/archlinux

#working with the main directory
WORKDIR /root
RUN mkdir /etc/shm
#adding it2s repository

#RUN echo '[it2s]' >> /etc/pacman.conf
#RUN echo 'SigLevel = Optional TrustAll' >> /etc/pacman.conf
#RUN echo 'Server = https= https://fpga.av.it.pt/repo/archlinux/$arch' >> /etc/pacman.conf

#installing requirements
RUN pacman -Syu
RUN pacman -S --noconfirm maven git jdk8-openjdk

ADD . /app
WORKDIR /app
RUN cp test_messages/* /etc/shm
RUN mvn clean install 
RUN cd demo-apps/cf-it2s-coap-server/
CMD mvn clean install