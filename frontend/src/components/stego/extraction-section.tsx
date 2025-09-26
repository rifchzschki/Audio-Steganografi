'use client';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type * as React from 'react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

type NullableFile = File | null;

export default function ExtractionSection() {
  const [stegoAudio, setStegoAudio] = useState<NullableFile>(null);
  const [key, setKey] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);

  const [extractedFile, setExtractedFile] = useState<NullableFile>(null);
  const [extractedName, setExtractedName] = useState('secret.bin');
  const [success, setSuccess] = useState<boolean | null>(null);

  function handleStegoInput(e: React.ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0] || null;
    setStegoAudio(f);
    setExtractedFile(null);
    setSuccess(null);
  }

  const canProcess = !!stegoAudio && !!key && key.length <= 25 && !isProcessing;

  async function onExtract() {
    if (!canProcess) return;
    setIsProcessing(true);
    setSuccess(null);
    setExtractedFile(null);

    await new Promise((r) => setTimeout(r, 1000));

    const suggested =
      stegoAudio?.name?.replace(/\.(mp3|mpeg)$/i, '') || 'secret';
    const outName = `${suggested}-extracted.bin`;
    setExtractedName(outName);
    const demoBlob = new Blob(
      [`[demo-secret]\nfrom=${stegoAudio?.name}\nkey=${key}\n`],
      {
        type: 'application/octet-stream',
      }
    );
    const out = new File([demoBlob], outName, {
      type: 'application/octet-stream',
    });
    setExtractedFile(out);
    setSuccess(true);

    toast('Extraction complete', {
      description: 'Secret message extracted. You can save it now.',
    });
    setIsProcessing(false);
  }

  function onDownloadSecret() {
    if (!extractedFile) return;
    const url = URL.createObjectURL(extractedFile);
    const a = document.createElement('a');
    a.href = url;
    a.download = extractedName || extractedFile.name;
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
  }

  useEffect(() => {
    if (success === false) {
      // TODO: handle failure case (not implemented in this mock)
    }
  }, [success]);

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Inputs & Uploads</CardTitle>
          <CardDescription>
            Provide the Stego-Audio file (MP3) and the stego key/seed used
            during insertion.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div className="grid gap-2">
            <Label htmlFor="stego-audio">Stego Audio File (MP3)</Label>
            <Input
              id="stego-audio"
              type="file"
              accept="audio/mpeg,audio/mp3"
              onChange={handleStegoInput}
            />
            <p className="text-sm text-muted-foreground">
              {stegoAudio
                ? `Selected: ${stegoAudio.name}`
                : 'Choose the MP3 that contains the hidden data.'}
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="extract-key">Stego Key/Seed (max 25 chars)</Label>
            <Input
              id="extract-key"
              value={key}
              maxLength={25}
              onChange={(e) => setKey(e.target.value)}
              placeholder="Enter your key"
            />
          </div>
        </CardContent>
        <CardFooter className="flex items-center justify-end gap-2">
          <Button
            onClick={onExtract}
            disabled={!canProcess}
            aria-busy={isProcessing}
          >
            {isProcessing ? 'Extracting...' : 'Extract Message'}
          </Button>
        </CardFooter>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Results & Download</CardTitle>
          <CardDescription>
            Confirmation and save options for the extracted secret message.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          {success === null && (
            <p className="text-sm text-muted-foreground">
              Extraction results will appear here.
            </p>
          )}
          {success === true && (
            <Alert>
              <AlertTitle>Extraction successful</AlertTitle>
              <AlertDescription>
                The secret message is ready to download.
              </AlertDescription>
            </Alert>
          )}
          {success === false && (
            <Alert variant="destructive">
              <AlertTitle>Extraction failed</AlertTitle>
              <AlertDescription>
                Could not recover the message. Check the key/seed and file.
              </AlertDescription>
            </Alert>
          )}

          <div className="grid gap-2 md:grid-cols-[1fr_auto] md:items-end">
            <div className="grid gap-2">
              <Label htmlFor="secret-filename">Save As</Label>
              <Input
                id="secret-filename"
                value={extractedName}
                onChange={(e) => setExtractedName(e.target.value)}
                placeholder="secret.bin"
              />
            </div>
            <Button onClick={onDownloadSecret} disabled={!extractedFile}>
              Save Secret Message
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
