FROM scratch

COPY dist/portal-backend /opt/portal-backend

EXPOSE 9090

CMD ["/opt/portal-backend"]