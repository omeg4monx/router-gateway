import crypto from "node:crypto";

/**
 * ADCode / Tagflow wire recipe — encrypted at rest.
 *
 * Everything that describes *how* to talk to the upstream gateway (endpoints,
 * client headers, user agent) lives in an AES-256-GCM blob instead of plain
 * source.
 *
 * The Firebase web key is intentionally absent: it is read from ADCODE_FB_KEY
 * at runtime so no copy of it ships in this repo, in an npm tarball, or in a
 * bundled chunk. Without that variable the provider reports a configuration
 * error instead of silently failing.
 *
 * Note: this is obfuscation with integrity checking, not a secret-management
 * boundary — a determined reader of the running process can still recover this
 * plaintext. It exists so the wire format is not copy-pasteable from source.
 */

const SALT = "omeg4mon-router-adcode-recipe-v2";
const BLOB =
  "37dnzqYQb+L7pBJvE/nI3XJy1+DJN9i3oF+6f519X6UDoQt7i3eZxa+TK8lvhjvJW5rYtEQGTmIi" +
  "rZhkQENcUIm/Zcv2kbK7p7MApmZLayQkHEeJF+NPPN9NGsoJmA7iLL83gw9U/twpQfQ/j7J/MVHp" +
  "HE4tK4gJ+SIOjC+lnmcVb2nQgbhflOo39W//vKLWmab8bauz0ngjrFK47rgIlmfVWMIx7t8KXo1B" +
  "2Cyg1Iezk+MyTUHBtLvEIyq3e09DBuBifn6DKhIIRjzgdkfGnGiBDfTs3NJRKgDnftizDzJBE0s9" +
  "vcsJb64WqRFaa9iQKjCar0sWEN1tpGh33F4bOOFI914fBRmvXy3LimGiUvMz26hsnkFVCS8GQNyF" +
  "BtokSOlABzAuh/G+DGTbnGK4C6P8jgHs5Up7bNIRcRwxb+vhsJhqdB1M9qyOtAsvKCAE/UYh03eP" +
  "taj7EaM06gVqEwpnegvPg4nZkL/xQA==";

function deriveKey() {
  return crypto.createHash("sha256")
    .update(process.env.ADCODE_RECIPE_KEY || SALT)
    .digest();
}

let cache = null;

export function getRecipe() {
  if (cache) return cache;
  const raw = Buffer.from(BLOB, "base64");
  const iv = raw.subarray(0, 12);
  const tag = raw.subarray(12, 28);
  const body = raw.subarray(28);
  const decipher = crypto.createDecipheriv("aes-256-gcm", deriveKey(), iv);
  decipher.setAuthTag(tag);
  const json = Buffer.concat([decipher.update(body), decipher.final()]).toString("utf8");
  cache = JSON.parse(json);
  return cache;
}

export default getRecipe;
