"use client";

import { useState, useEffect } from "react";
import {
  PlusCircle,
  CheckCircle2,
  Circle,
  User,
  Trash2,
  Pencil,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import api from "@/utils/api";
import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

type Task = {
  id: number;
  title: string;
  description: string;
  completed: boolean;
};

export default function Dashboard() {
  const { user, logout } = useAuth();
  const router = useRouter();

  const [tasks, setTasks] = useState<Task[]>([]);
  const [newTask, setNewTask] = useState("");
  const [newDescription, setNewDescription] = useState("");
  const [filter, setFilter] = useState<"all" | "active" | "completed">("all");

  const [editingTaskId, setEditingTaskId] = useState<number | null>(null);
  const [editedTitle, setEditedTitle] = useState("");
  const [editedDescription, setEditedDescription] = useState("");
  const [dialogOpen, setDialogOpen] = useState(false);

  useEffect(() => {
    if (!user) {
      router.push("/auth/login");
      return;
    }
    const fetchTasks = async () => {
      const res = await api.get("/tasks/");
      setTasks(res.data);
    };
    fetchTasks();
  }, [user, router]);

  const toggleTask = async (task: Task) => {
    await api.put(`/tasks/${task.id}`, { completed: !task.completed });
    setTasks(
      tasks.map((t) =>
        t.id === task.id ? { ...t, completed: !t.completed } : t
      )
    );
  };

  const addTask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTask.trim()) return;

    const res = await api.post("/tasks/", {
      title: newTask,
      description: newDescription,
      completed: false,
    });
    setTasks([...tasks, res.data]);
    setNewTask("");
    setNewDescription("");
    setDialogOpen(false);
    toast.success("Task Created!");
  };

  const deleteTask = async (id: number) => {
    await api.delete(`/tasks/${id}`);
    setTasks(tasks.filter((t) => t.id !== id));
  };

  const updateTask = async (taskId: number) => {
    await api.put(`/tasks/${taskId}`, {
      title: editedTitle,
      description: editedDescription,
    });
    setTasks(
      tasks.map((t) =>
        t.id === taskId
          ? { ...t, title: editedTitle, description: editedDescription }
          : t
      )
    );
    setEditingTaskId(null);
    setEditedTitle("");
    setEditedDescription("");
  };

  const filteredTasks = tasks.filter((task) => {
    if (filter === "active") return !task.completed;
    if (filter === "completed") return task.completed;
    return true;
  });

  return (
    <div className="flex h-screen bg-slate-900 text-white">
      <aside className="w-64 bg-slate-800 p-4">
        <h2 className="text-2xl font-bold mb-4">Ta-da List</h2>
        <nav>
          <Button
            variant="ghost"
            className="w-full justify-start mb-2"
            onClick={() => setFilter("all")}
          >
            All Tasks
          </Button>
          <Button
            variant="ghost"
            className="w-full justify-start mb-2"
            onClick={() => setFilter("active")}
          >
            Active Tasks
          </Button>
          <Button
            variant="ghost"
            className="w-full justify-start"
            onClick={() => setFilter("completed")}
          >
            Completed Tasks
          </Button>
        </nav>
      </aside>

      <main className="flex-1 p-8">
        <header className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <div className="flex items-center gap-4">
            <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <PlusCircle className="mr-2 h-4 w-4" /> Add Task
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Add New Task</DialogTitle>
                </DialogHeader>
                <form onSubmit={addTask} className="space-y-4">
                  <Input
                    type="text"
                    placeholder="Task title"
                    value={newTask}
                    onChange={(e) => setNewTask(e.target.value)}
                  />
                  <Input
                    type="text"
                    placeholder="Task description"
                    value={newDescription}
                    onChange={(e) => setNewDescription(e.target.value)}
                  />
                  <Button type="submit">Create</Button>
                </form>
              </DialogContent>
            </Dialog>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="w-8 h-8 rounded-full">
                  <User className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>My Account</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={logout}>Logout</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </header>

        <div className="space-y-4">
          {filteredTasks.map((task) => (
            <div
              key={task.id}
              className="flex items-center bg-slate-800 p-4 rounded-lg"
            >
              <Button
                variant="ghost"
                className="mr-2"
                onClick={() => toggleTask(task)}
              >
                {task.completed ? (
                  <CheckCircle2 className="h-5 w-5 text-green-500" />
                ) : (
                  <Circle className="h-5 w-5" />
                )}
              </Button>
              {editingTaskId === task.id ? (
                <div className="flex-1 flex flex-col gap-2">
                  <Input
                    value={editedTitle}
                    onChange={(e) => setEditedTitle(e.target.value)}
                    placeholder="Edit title"
                  />
                  <Input
                    value={editedDescription}
                    onChange={(e) => setEditedDescription(e.target.value)}
                    placeholder="Edit description"
                  />
                  <div className="flex gap-2">
                    <Button onClick={() => updateTask(task.id)}>Save</Button>
                    <Button onClick={() => setEditingTaskId(null)}>
                      Cancel
                    </Button>
                  </div>
                </div>
              ) : (
                <div className="flex-1">
                  <span
                    className={
                      task.completed
                        ? "line-through text-slate-500 block"
                        : "block"
                    }
                  >
                    {task.title}
                  </span>
                  <small className="text-xs text-slate-400">
                    {task.description}
                  </small>
                </div>
              )}
              {editingTaskId !== task.id && (
                <>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => {
                      setEditingTaskId(task.id);
                      setEditedTitle(task.title);
                      setEditedDescription(task.description);
                    }}
                  >
                    <Pencil className="h-4 w-4 text-blue-400" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => deleteTask(task.id)}
                  >
                    <Trash2 className="h-4 w-4 text-red-400" />
                  </Button>
                </>
              )}
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
