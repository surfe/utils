package urls

import (
	"strings"
)

// IsPublicDomain checks if domain starts with known public (social media, website builder) domains
func IsPublicDomain(domain string) bool {
	var pd []string
	pd = append(pd, websiteBuilders...)
	pd = append(pd, socialMediaDomains...)

	for _, publicDomain := range pd {
		if strings.EqualFold(domain, publicDomain) {
			return true
		}
	}

	return false
}

// IsURLShortenerDomain checks if domain is a known url shortener domain
func IsURLShortenerDomain(domain string) bool {
	return urlShortenerDomains[domain]
}

var websiteBuilders = []string{
	"zohosites.com",
	"weebly.com",
	"governor.io",
	"sitebuilder.com",
	"blogger.com",
	"jimbo.com",
	"site123.com",
	"doodlekit.com",
	"wordpress.com",
	"wix.com",
	"wix.net",
	"squarespace.com",
	"godaddy.com",
}

var socialMediaDomains = []string{
	"facebook.com",
	"twitter.com",
	"instagram.com",
	"whatsapp.com",
	"tiktok.com",
	"reddit.com",
	"linkedin.com",
	"linktr.ee",
	"vk.com",
	"discord.com",
	"pinterest.com",
	"ok.ru",
	"zhihu.com",
	"messenger.com",
	"line.me",
	"telegram.org",
	"tumblr.com",
	"namu.wiki",
	"nextdoor.com",
	"ameblo.jp",
	"weibo.com",
	"ppgames.net",
	"redd.it",
	"slack.com",
	"zalo.me",
	"patreon.com",
	"livejournal.com",
	"slideshare.net",
	"snapchat.com",
	"discordapp.com",
	"hatenablog.com",
	"hczog.com",
	"omegle.com",
	"fb.com",
	"pinterest.es",
	"snaptik.app",
	"ssstik.io",
	"gotrackier.com",
	"bakusai.com",
	"pinterest.com.mx",
	"51dongshi.com",
	"ptt.cc",
	"fb.watch",
	"pinterest.co.uk",
	"kwai.com",
	"pinterest.fr",
	"ninisite.com",
	"bp.blogspot.com",
	"dcard.tw",
	"youtubekids.com",
	"ameba.jp",
}
