export async function Decrypt(ciphertextB64: string, keyB64: string) {
  const rawData = Uint8Array.from(atob(ciphertextB64), (c) => c.charCodeAt(0));
  const keyBytes = Uint8Array.from(atob(btoa(keyB64)), (c) => c.charCodeAt(0));

  const nonce = rawData.slice(0, 12);
  const ciphertext = rawData.slice(12);

  const cryptoKey = await crypto.subtle.importKey(
    "raw",
    keyBytes,
    { name: "AES-GCM" },
    false,
    ["decrypt"],
  );

  const decrypted = await crypto.subtle.decrypt(
    { name: "AES-GCM", iv: nonce },
    cryptoKey,
    ciphertext,
  );

  return new TextDecoder().decode(decrypted);
}
