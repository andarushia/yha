{ config, lib, pkgs, ... }:

{
  config.project.name = "anyapg";
  config.services.postgres = {
    service.image = "postgres:16";
    service.volumes = [ "${toString ./.}/postgres-data:/var/lib/postgresql/data" ];
    service.environment.POSTGRES_PASSWORD = "sqlxpass";
    service.environment.POSTGRES_USER = "anya";
    service.environment.POSTGRES_DB = "anyadb";
  };
}
