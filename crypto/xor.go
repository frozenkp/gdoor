package crypto

func XOR(text, key string)string{
  output := ""
  for i, _ := range text {
    output += string(text[i] ^ key[i % len(key)])
  }
  return output
}
