#!/usr/bin/env bash

SESSION="ecommerce"
ROOT="$HOME/Dev/go/ecommerce/services"

# Don't create it again if it already exists
if tmux has-session -t "$SESSION" 2>/dev/null; then
  echo "Session '$SESSION' already exists."
  tmux attach -t "$SESSION"
  exit 0
fi

# Create first window
tmux new-session -d -s "$SESSION" -n gateway

tmux send-keys -t "$SESSION:gateway" \
  "cd $ROOT/api-gateway && air" C-m

# Auth
tmux new-window -t "$SESSION" -n auth
tmux send-keys -t "$SESSION:auth" \
  "cd $ROOT/auth && air" C-m

# Product
tmux new-window -t "$SESSION" -n product
tmux send-keys -t "$SESSION:product" \
  "cd $ROOT/product && air" C-m

# Inventory
tmux new-window -t "$SESSION" -n inventory
tmux send-keys -t "$SESSION:inventory" \
  "cd $ROOT/inventory && air" C-m

# Order
tmux new-window -t "$SESSION" -n order
tmux send-keys -t "$SESSION:order" \
  "cd $ROOT/order && air" C-m

# Payment
tmux new-window -t "$SESSION" -n payment
tmux send-keys -t "$SESSION:payment" \
  "cd $ROOT/payment && air" C-m

# Shipping
tmux new-window -t "$SESSION" -n shipping
tmux send-keys -t "$SESSION:shipping" \
  "cd $ROOT/shipping && air" C-m

# Notification
tmux new-window -t "$SESSION" -n Notification
tmux send-keys -t "$SESSION:Notification" \
  "cd $ROOT/notification && air" C-m

# Select the first window
tmux select-window -t "$SESSION:gateway"

# Attach
tmux attach -t "$SESSION"
