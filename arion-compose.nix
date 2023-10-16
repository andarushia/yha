{ config, lib, pkgs, ... }:

{
  config.project.name = "anyapg";
  config.services.postgres = {
    service.image = "postgres:16";
    service.volumes = [ "${toString ./.}/postgres-data:/var/lib/postgresql/data" ];
    service.environment = {
      POSTGRES_PASSWORD = "sqlxpass";
      POSTGRES_USER = "anya";
      POSTGRES_DB = "anyadb";
    };
    service.ports = [ "5432:5432" ];
  };
  config.services.pgadmin4 = {
    service.image = "dpage/pgadmin4:7";
    service.environment = {
      PGADMIN_DEFAULT_EMAIL = "andalusia@porque.org";
      PGADMIN_DEFAULT_PASSWORD = "sqlxpass";
      PGADMIN_LISTEN_PORT = 5050;
    };
    service.ports = [ "5050:5050" ];
  };
}
