package command

import (
	"context"
	"fmt"
	"os"

	"github.com/darkseear/gophkeeper/client/internal/client"
	"github.com/darkseear/gophkeeper/client/internal/storage"
	"github.com/spf13/cobra"
)

// NewRootCmd - рут команды.
func NewRootCmd(version, buildDate string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "gophkeeper",
		Short:   "Secure password manager",
		Version: fmt.Sprintf("%s (built at %s)", version, buildDate),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initClientInContext(cmd)
		},
	}

	// Добавляем флаги
	rootCmd.PersistentFlags().StringP("server", "s", "localhost:50051", "Server address")
	rootCmd.PersistentFlags().StringP("password", "p", "master", "Master password")
	rootCmd.PersistentFlags().StringP("bdname", "b", "gophkeeper.db", "Name local bd")

	// Добавляем подкоманды
	rootCmd.AddCommand(
		newRegisterCmd(),
		newLoginCmd(),
		newSecretCmd(),
		newGetCmd(),
	)

	return rootCmd
}

// newRegisterCmd - команда регистрации.
func newRegisterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "register [login] [password]",
		Short: "Register new user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Получаем клиент из контекста
			cli, err := getClient(cmd)
			if err != nil {
				fmt.Printf("Client error: %v\n", err)
				return
			}

			// Выполняем регистрацию
			if err := cli.Register(cmd.Context(), args[0], args[1]); err != nil {
				fmt.Printf("Registration failed: %v\n", err)
				return
			}
			fmt.Println("Successfully registered")
		},
	}
}

// newLoginCmd - команда для авторизации юзеров.
func newLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login [login] [password]",
		Short: "Authenticate user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Получаем клиент из контекста
			cli, err := getClient(cmd)
			if err != nil {
				fmt.Printf("Client error: %v\n", err)
				return
			}

			// Выполняем регистрацию
			if err := cli.Login(cmd.Context(), args[0], args[1]); err != nil {
				fmt.Printf("Login failed: %v\n", err)
				return
			}
			fmt.Println("Successfully logged in")
		},
	}
}

// newSecretCmd - команда для добаления секретов.
func newSecretCmd() *cobra.Command {
	var (
		secretType  string
		description string
		secretFile  string
	)

	secretCmd := &cobra.Command{

		Use:   "add-secret [data]",
		Short: "Add new secret (for text: [data], for card: [number] [expiry] [cvv], for binary: use --file)",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {

			// Чтение данных в зависимости от типа
			var data []byte
			var err error

			cli, err := getClient(cmd)
			if err != nil {
				fmt.Printf("Client error: %v\n", err)
				return
			}

			switch secretType {
			case "text":
				if len(args) < 1 {
					fmt.Println("Text secret requires at least 1 argument")
					return
				}
				fmt.Println(args[0])
				data = []byte(args[0])
			case "binary":
				if secretFile == "" {
					fmt.Println("Binary secret requires --file flag")
					return
				}
				data, err = os.ReadFile(secretFile)
				if err != nil {
					fmt.Printf("Error reading file: %v\n", err)
					return
				}
			case "card":
				if len(args) < 3 {
					fmt.Println("Card secret requires 3 arguments: number, expiry, cvv")
					return
				}
				data = make([]byte, 0, 128)
				data = fmt.Appendf(data, `{"number":"%s","expiry":"%s","cvv":"%s"}`, args[0], args[1], args[2])
			default:
				fmt.Println("Invalid secret type")
				return
			}

			if err := cli.AddSecret(cmd.Context(), secretType, description, data); err != nil {
				fmt.Printf("Error adding secret: %v\n", err)
				return
			}

			fmt.Println("Secret added successfully")
		},
	}

	// Флаги команды
	secretCmd.Flags().StringVarP(&secretType, "type", "t", "", "Secret type (text|binary|card)")
	secretCmd.Flags().StringVarP(&description, "desc", "d", "", "Secret description")
	secretCmd.Flags().StringVarP(&secretFile, "file", "f", "", "File path for binary data")

	return secretCmd
}

// newGetCmd - команда получения приватных данных.
func newGetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get-secret [secret-id]",
		Short: "Get secret data by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := getClient(cmd)
			if err != nil {
				fmt.Printf("Client error: %v\n", err)
				return
			}

			// Получаем секрет
			secret, err := cli.GetSecret(cmd.Context(), args[0])
			if err != nil {
				fmt.Printf("Error getting secret: %v\n", err)
				os.Exit(1)
			}

			// Выводим результат в зависимости от типа
			switch secret.Type {
			case "text":
				fmt.Printf("Text content: %s\n", string(secret.Data))
			case "card":
				fmt.Printf("Card: %s",
					string(secret.Data),
				)
			default:
				fmt.Printf("Secret ID: %s\nType: %s\nMetadata: %v\nData: [%d bytes]\n",
					secret.ID, secret.Type, secret.Metadata, len(secret.Data))
			}
		},
	}
	return getCmd
}

// initClientInContext - создает клиента и сохраняет его в контекст команды.
// Используйте перед запуском команд, требующих клиента.

func initClientInContext(cmd *cobra.Command) error {
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		return fmt.Errorf("failed to get server flag: %w", err)
	}
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return fmt.Errorf("failed to get password flag: %w", err)
	}
	bdname, err := cmd.Flags().GetString("bdname")
	if err != nil {
		return fmt.Errorf("failed to get bdname flag: %w", err)
	}

	network, err := client.NewNetworkClient(server)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	local, err := storage.NewLocalStorage(bdname)
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	cli, err := client.NewGophkeeperClient(network, local, password)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	ctx := cmd.Context()
	ctx = context.WithValue(ctx, "gophkeeper_client", cli)
	cmd.SetContext(ctx)
	return nil
}

// getClient - Вспомогательная функция для получения клиента из контекста.
func getClient(cmd *cobra.Command) (*client.GophKeeperClient, error) {
	ctx := cmd.Context()
	cli, ok := ctx.Value("gophkeeper_client").(*client.GophKeeperClient)
	if !ok || cli == nil {
		return nil, fmt.Errorf("client not found in context")
	}
	return cli, nil
}
