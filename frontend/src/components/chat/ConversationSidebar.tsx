'use client';

import { useState, useEffect } from 'react';
import { useChatStore } from '@/store/chat';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Plus, Settings } from 'lucide-react';
import { ThemeToggle } from '@/components/theme-toggle';

export default function ConversationSidebar() {
  const {
    conversations,
    currentConversation,
    models,
    isLoading,
    loadConversations,
    loadModels,
    createConversation,
    selectConversation,
    deleteConversation,
  } = useChatStore();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [newTitle, setNewTitle] = useState('');
  const [selectedModel, setSelectedModel] = useState('');

  useEffect(() => {
    loadConversations();
    loadModels();
  }, [loadConversations, loadModels]);

  const handleCreateConversation = async () => {
    if (newTitle.trim() && selectedModel) {
      try {
        await createConversation(newTitle.trim(), selectedModel);
        setNewTitle('');
        setSelectedModel('');
        setIsCreateOpen(false);
      } catch (error) {
        console.error('Failed to create conversation:', error);
      }
    }
  };

  const handleSelectConversation = (id: string) => {
    selectConversation(id);
  };

  const handleDeleteConversation = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (confirm('Are you sure you want to delete this conversation?')) {
      try {
        await deleteConversation(id);
      } catch (error) {
        console.error('Failed to delete conversation:', error);
      }
    }
  };

  return (
    <div className="w-80 border-r bg-muted/50 flex flex-col">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">Conversations</CardTitle>
          <div className="flex items-center gap-2">
            <ThemeToggle />
            <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
              <DialogTrigger asChild>
                <Button size="sm" variant="outline">
                  <Plus className="h-4 w-4 mr-2" />
                  New
                </Button>
              </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Create New Conversation</DialogTitle>
              </DialogHeader>
              <div className="space-y-4">
                <div>
                  <Label htmlFor="title">Title</Label>
                  <Input
                    id="title"
                    value={newTitle}
                    onChange={(e) => setNewTitle(e.target.value)}
                    placeholder="Enter conversation title"
                  />
                </div>
                <div>
                  <Label htmlFor="model">Model</Label>
                  <Select value={selectedModel} onValueChange={setSelectedModel}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select a model" />
                    </SelectTrigger>
                    <SelectContent>
                      {models.map((model) => (
                        <SelectItem key={model.id} value={model.id}>
                          {model.name} ({model.provider})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex justify-end gap-2">
                  <Button variant="outline" onClick={() => setIsCreateOpen(false)}>
                    Cancel
                  </Button>
                  <Button 
                    onClick={handleCreateConversation}
                    disabled={!newTitle.trim() || !selectedModel}
                  >
                    Create
                  </Button>
                </div>
              </div>
            </DialogContent>
          </Dialog>
          </div>
        </div>
      </CardHeader>

      <CardContent className="flex-1 overflow-y-auto space-y-2">
        {isLoading && (!conversations || conversations.length === 0) ? (
          <div className="text-center text-muted-foreground py-8">
            Loading conversations...
          </div>
        ) : (!conversations || conversations.length === 0) ? (
          <div className="text-center text-muted-foreground py-8">
            <p>No conversations yet</p>
            <p className="text-sm">Create your first conversation to get started.</p>
          </div>
        ) : (
          conversations?.map((conversation) => (
            <Card
              key={conversation.id}
              className={`cursor-pointer transition-colors hover:bg-accent ${
                currentConversation?.id === conversation.id ? 'bg-accent' : ''
              }`}
              onClick={() => handleSelectConversation(conversation.id)}
            >
              <CardContent className="p-3">
                <div className="flex items-center justify-between">
                  <div className="flex-1 min-w-0">
                    <h4 className="font-medium text-sm truncate">
                      {conversation.title}
                    </h4>
                    <p className="text-xs text-muted-foreground">
                      {conversation.model}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {new Date(conversation.updated_at).toLocaleDateString()}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0"
                    onClick={(e) => handleDeleteConversation(conversation.id, e)}
                  >
                    <Settings className="h-4 w-4" />
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </CardContent>
    </div>
  );
}