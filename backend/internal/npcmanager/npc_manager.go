package npcmanager

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nancy010006/MakeMoney/internal/npc"
)

type NPCManager struct {
	npcs map[string]*npc.NPC
	mu   sync.Mutex
	db   *sql.DB
}

func New(db *sql.DB) *NPCManager {
	return &NPCManager{
		npcs: make(map[string]*npc.NPC),
		db:   db,
	}
}

func (nm *NPCManager) StartNPCHandler(c *gin.Context) {
	npcID := c.Query("id")
	if npcID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NPC ID is required"})
		return
	}

	nm.mu.Lock()
	defer nm.mu.Unlock()

	if _, exists := nm.npcs[npcID]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "NPC already running"})
		return
	}

	newNPC := npc.New(npcID, nm.db)
	go newNPC.Start()
	nm.npcs[npcID] = newNPC

	c.JSON(http.StatusOK, gin.H{"status": "NPC started"})
}

func (nm *NPCManager) StopNPCHandler(c *gin.Context) {
	npcID := c.Query("id")
	if npcID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NPC ID is required"})
		return
	}

	nm.mu.Lock()
	defer nm.mu.Unlock()

	if npc, exists := nm.npcs[npcID]; exists {
		npc.Stop()
		delete(nm.npcs, npcID)
		c.JSON(http.StatusOK, gin.H{"status": "NPC stopped"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "NPC not found"})
	}
}

func (nm *NPCManager) ListNPCsHandler(c *gin.Context) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	npcList := make([]string, 0, len(nm.npcs))
	for id := range nm.npcs {
		npcList = append(npcList, id)
	}

	c.JSON(http.StatusOK, gin.H{"npcs": npcList})
}

func (nm *NPCManager) StopAll() {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	for _, npc := range nm.npcs {
		npc.Stop()
	}
	nm.npcs = make(map[string]*npc.NPC)
}
