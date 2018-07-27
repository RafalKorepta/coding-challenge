FROM scratch

COPY dist/portal-backend /opt/portal-backend

EXPOSE 9091

CMD ["/opt/portal-backend"]