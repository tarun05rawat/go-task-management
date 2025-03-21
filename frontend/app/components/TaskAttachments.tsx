"use client";

import { useEffect, useState } from "react";
import Image from "next/image";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import api from "@/utils/api";

export default function TaskAttachments({ taskId }: { taskId: number }) {
  const [open, setOpen] = useState(false);
  const [attachments, setAttachments] = useState<string[]>([]);
  const [selectedFiles, setSelectedFiles] = useState<FileList | null>(null);

  // Fetch existing attachments on dialog open
  useEffect(() => {
    if (open) {
      api.get(`/tasks/${taskId}/attachments`).then((res) => {
        setAttachments(res.data.attachments);
      });
    }
  }, [open, taskId]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedFiles(e.target.files);
  };

  const handleUpload = async () => {
    if (!selectedFiles) return;

    const formData = new FormData();
    Array.from(selectedFiles).forEach((file) => {
      formData.append("files", file);
    });

    await api.post(`/tasks/${taskId}/upload`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });

    // Refresh attachments list after upload
    const res = await api.get(`/tasks/${taskId}/attachments`);
    setAttachments(res.data.attachments);
    setSelectedFiles(null);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline">Manage Attachments</Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Attachments for Task #{taskId}</DialogTitle>
        </DialogHeader>

        {/* Upload Section */}
        <div className="space-y-2">
          <Input type="file" multiple onChange={handleFileChange} />
          <Button onClick={handleUpload} disabled={!selectedFiles}>
            Upload Files
          </Button>
        </div>

        {/* Attachments Preview */}
        <div className="grid grid-cols-2 gap-2 mt-4">
          {attachments.map((url, idx) => (
            <div key={idx} className="border rounded p-2 bg-slate-800">
              {url.endsWith(".jpg") || url.endsWith(".png") ? (
                <Image
                  src={url}
                  alt="attachment"
                  width={200}
                  height={100}
                  className="w-full h-24 object-cover rounded"
                />
              ) : (
                <a
                  href={url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-blue-300 underline"
                >
                  {url.split("/").pop()}
                </a>
              )}
            </div>
          ))}
        </div>
      </DialogContent>
    </Dialog>
  );
}
