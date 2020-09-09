File.open('internal/app/web/testdata/link.dat', 'w') do |f|
  s = [
    1, 2, 4, 4, 0, 3, 0,
    3, 0, 1, 4, 1, 2,
  ].pack("i*")
  f.write s
end
