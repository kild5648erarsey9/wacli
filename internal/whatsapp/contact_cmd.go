package whatsapp

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewContactCmd returns a cobra command for sending a contact message.
func NewContactCmd() *cobra.Command {
	var (
		to        string
		name      string
		phone     string
		phoneType string
	)

	cmd := &cobra.Command{
		Use:   "contact",
		Short: "Send a contact card via WhatsApp",
		RunE: func(cmd *cobra.Command, args []string) error {
			phoneID := os.Getenv("WA_PHONE_ID")
			token := os.Getenv("WA_ACCESS_TOKEN")

			client, err := NewClient(phoneID, token)
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}

			contact := Contact{
				Name: ContactName{FormattedName: name},
			}
			if phone != "" {
				contact.Phones = []ContactPhone{{Phone: phone, Type: phoneType}}
			}

			if err := client.SendContactMessage(context.Background(), to, []Contact{contact}); err != nil {
				return fmt.Errorf("send contact: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Contact sent to %s\n", to)
			return nil
		},
	}

	cmd.Flags().StringVarP(&to, "to", "t", "", "Recipient phone number (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Contact formatted name (required)")
	cmd.Flags().StringVarP(&phone, "phone", "p", "", "Contact phone number")
	cmd.Flags().StringVar(&phoneType, "phone-type", "MOBILE", "Contact phone type (MOBILE, HOME, WORK)")
	_ = cmd.MarkFlagRequired("to")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
