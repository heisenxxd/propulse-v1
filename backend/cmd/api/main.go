package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"propulse/handler"
	"propulse/shared/logger"
)

func main() {
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	
	logger.Info("Aguardando banco de dados iniciar")
	timer := time.NewTimer(5 * time.Second)
    <-timer.C
	
	logger.Info("Configuração do banco de dados iniciada...")
	dbpool, err := DbConfig(ctx)
	if err != nil {
		logger.Error("Erro fatal ao configurar o banco de dados", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	logger.Info("Conexão com o banco de dados estabelecida com sucesso!")
	
	router := gin.Default()

	port := os.Getenv("PORT"); if port == "" {
    port = "8080"
    }

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	logger.Info("Inicializando as rotas dos serviços!")
	handler.SetupServices(dbpool, router)

	go SetupServer(server)

	<-ctx.Done()

	logger.Info("Sinal de desligamento recebido. Iniciando graceful shutdown...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Erro no graceful shutdown do servidor", err)
	}

	logger.Info("Servidor finalizado com sucesso.")
}

func SetupServer(server *http.Server) {
	logger.Info("Servidor iniciado", zap.String("porta", server.Addr))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Erro ao iniciar o servidor HTTP", err)
		os.Exit(1)
	}
}

func DbConfig(ctx context.Context) (*pgxpool.Pool, error) {
	connectString := os.Getenv("DATABASE_URL")
	if connectString == "" {
		return nil, fmt.Errorf("DATABASE_URL não está definido no ambiente")
	}

	config, err := pgxpool.ParseConfig(connectString)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear a string de conexão: %w", err)
	}

	connCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(connCtx, config)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	if err := dbpool.Ping(connCtx); err != nil {
		return nil, fmt.Errorf("erro ao pingar o banco de dados: %w", err)
	}

	return dbpool, nil
}