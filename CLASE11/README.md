# Clase 11

## Instalación de Go

Actualizar e instalar Go en sistemas basados en Debian/Ubuntu:

```bash
sudo apt update
sudo apt upgrade
sudo apt install golang-go
```

## Configuración de Permisos para AWS S3

Política de permisos para hacer un bucket S3 accesible públicamente:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::NOMBRE-DEL-BUCKET/*"
    }
  ]
}
```

> **Nota**: Reemplazar `NOMBRE-DEL-BUCKET` con el nombre de tu bucket S3.
