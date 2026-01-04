<xsl:for-each select="//player">
<script>
window.dps = <xsl:value-of select="dps/mean"/>;
</script>
</xsl:for-each>