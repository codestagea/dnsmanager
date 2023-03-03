// from vben
import { encrypt, decrypt } from 'crypto-js/aes';
import { parse } from 'crypto-js/enc-utf8';
import pkcs7 from 'crypto-js/pad-pkcs7';
import ECB from 'crypto-js/mode-ecb';
import md5 from 'crypto-js/md5';
import UTF8 from 'crypto-js/enc-utf8';
import Base64 from 'crypto-js/enc-base64';
import { JSEncrypt } from 'jsencrypt'

export interface EncryptionParams {
  key: string;
}

export class AesEncryption {
  private key;

  constructor(opt: Partial<EncryptionParams> = {}) {
    const { key } = opt;
    if (key) {
      this.key = parse(key);
    }
  }

  get getOptions() {
    return {
      mode: ECB,
      padding: pkcs7,
    };
  }

  encrypt(cipherText: string) {
    return encrypt(cipherText, this.key, this.getOptions).toString();
  }

  decrypt(cipherText: string) {
    return decrypt(cipherText, this.key, this.getOptions).toString(UTF8);
  }
}

export class RsaEncryption {
  private key;

  constructor(key: string) {
    this.key = decodeByBase64(key)
  }

  encrypt(cipherText: string) {
    let encryptor = new JSEncrypt();
    encryptor.setPublicKey(this.key);
    return encryptor.encrypt(cipherText);
  }
}

export function encryptByBase64(cipherText: string) {
  return UTF8.parse(cipherText).toString(Base64);
}

export function decodeByBase64(cipherText: string) {
  return Base64.parse(cipherText).toString(UTF8);
}

export function encryptByMd5(password: string) {
  return md5(password).toString();
}

const rsa_pub_key = import.meta.env.VITE_RSA_PUBLIC_KEY

export const passwdEnc = new RsaEncryption(rsa_pub_key)
