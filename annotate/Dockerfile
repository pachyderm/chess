FROM ubuntu

# Install Stockfish.
RUN apt-get update && apt-get install -y stockfish curl && rm -rf /var/lib/apt/lists/*

# Add the annotate binary.
ADD annotate /
