require 'net/http'

class DroplrGraber

  def run
    chars = ('A'..'Z').to_a + ('a'..'z').to_a + ('0'..'9').to_a
    count = 0

    chars.each do |first_letter|
      chars.each do |second_letter|
        chars.each do |third_letter|
          chars.each do |fourth_letter|
            count += 1
            link = "AA#{third_letter}#{fourth_letter}"
            link = "#{first_letter}#{second_letter}#{third_letter}#{fourth_letter}"
            # p link
            # break if count == 50
            p "#{count}. #{link}: #{grab(link)}"
          end
        end
      end
    end
  end

  def grab(link)
    response = Net::HTTP.get_response(URI("http://d.pr/i/#{link}"))
    response.code == '200'
  end

end

grabber = DroplrGraber.new
grabber.run
