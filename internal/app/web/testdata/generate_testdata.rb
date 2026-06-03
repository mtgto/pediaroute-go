File.open('internal/app/web/testdata/forward_link.dat', 'w') do |f|
  s = [1, 2, 4, 4, 0, 3, 0].pack("I*")
  f.write s
end

File.open('internal/app/web/testdata/backward_link.dat', 'w') do |f|
  s = [3, 0, 1, 4, 1, 2].pack("I*")
  f.write s
end
