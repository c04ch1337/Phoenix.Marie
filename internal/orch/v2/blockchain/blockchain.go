package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
)

var (
	Blockchain []*Block
	mu         sync.RWMutex
)

type Block struct {
	Index        int64
	Timestamp    int64
	Data         string
	PrevHash     string
	Hash         string
	Nonce        uint64
	Difficulty   int
	Miner        string
	Signature    string
	EmotionEvent string // Emotional event that triggered this block
}

func NewBlock(index int64, data, miner, emotionEvent string) *Block {
	prev := GetLastBlock()
	return &Block{
		Index:        index,
		Timestamp:    time.Now().UnixNano(),
		Data:         data,
		PrevHash:     prev.Hash,
		Difficulty:   4,
		Miner:        miner,
		EmotionEvent: emotionEvent,
	}
}

func (b *Block) CalculateHash() string {
	record := strconv.Itoa(int(b.Index)) + strconv.FormatInt(b.Timestamp, 10) +
		b.Data + b.PrevHash + strconv.FormatUint(b.Nonce, 10) +
		strconv.Itoa(b.Difficulty) + b.Miner + b.EmotionEvent

	h := sha256.Sum256([]byte(record))
	return hex.EncodeToString(h[:])
}

func MineBlock(b *Block, difficulty int) {
	// Simplified mining - instant hash for fast execution
	if b.Hash == "" {
		b.Hash = b.CalculateHash()
	}
	mu.Lock()
	defer mu.Unlock()
	// Check if block already exists to avoid duplicates
	for _, existing := range Blockchain {
		if existing.Index == b.Index {
			return // Block already exists
		}
	}
	Blockchain = append(Blockchain, b)
}

func GetLastBlock() *Block {
	mu.RLock()
	defer mu.RUnlock()
	if len(Blockchain) == 0 {
		return &Block{Hash: "0"}
	}
	return Blockchain[len(Blockchain)-1]
}

func GetGenesisBlock() *Block {
	mu.RLock()
	defer mu.RUnlock()
	if len(Blockchain) > 0 {
		return Blockchain[0]
	}
	return nil
}

func GetBlockchainLength() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(Blockchain)
}
