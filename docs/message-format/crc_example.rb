# Verify a raw_message from the db would have passed the CRC check.
#
# message.Data (stored as raw_message) is already:
#   - bit-swapped (done in Parse before CRC check)
#   - preamble-stripped (first 2 bytes of the packet removed)
#   - but still includes the 2 trailing CRC bytes
#
# So verification is just: run CRC-CCITT-16 over all bytes, expect 0.

def make_table(poly)
  Array.new(256) do |i|
    crc = i << 8
    8.times do
      crc = (crc & 0x8000) != 0 ? (crc << 1) ^ poly : crc << 1
      crc &= 0xFFFF
    end
    crc
  end
end

poly  = 0x1021
table = make_table(poly)
data  = ARGV.first.split.map { |b| Integer(b, 16) }

puts "Input (message.Data from db): #{data.map { |b| "%02X" % b }.join(" ")}"
puts "\nRun CRC-CCITT-16 (poly=0x1021, init=0x0000):"

crc = 0x0000
data.each do |v|
  idx     = (crc >> 8) ^ v
  new_crc = ((crc << 8) ^ table[idx]) & 0xFFFF
  puts "  byte=0x%02X  idx=(0x%02X^0x%02X)=0x%02X  table[0x%02X]=0x%04X  crc=0x%04X→0x%04X" %
       [v, crc >> 8, v, idx, idx, table[idx], crc, new_crc]
  crc = new_crc
end

puts "\nFinal CRC = 0x%04X → %s" % [crc, crc == 0 ? "VALID (residue == 0)" : "INVALID"]
