touch backend/{Dockerfile,.dockerignore}
cd backend/ && go mod init SlotSwapper

go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/gin-contrib/cors
go get -u github.com/joho/godotenv
go get -u github.com/rs/zerolog
go get -u github.com/ulule/limiter/v3

touch main.go
mkdir -p config models dto utils middleware services controllers routes tests

touch config/database.go models/{user,event,swap_request}.go dto/{auth,event,swap}_tdo.go utils/{transaction,jwt,errors}.go middleware/{auth,cors,rate_limit,logger,error_handler}.go services/swap_service.go controllers/{auth,event,swap,health}_controller.go routes/route.go

touch .env

sudo -u postgres createuser --superuser samarth
createdb -U postgres slotswapper_dev
ALTER SCHEMA public OWNER TO slotswapper_user;
GRANT USAGE ON SCHEMA public TO slotswapper_user;
GRANT CREATE ON SCHEMA public TO slotswapper_user;
